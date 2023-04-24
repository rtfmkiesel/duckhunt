package user32

import (
	"duckhunt/pkg/config"
	"duckhunt/pkg/logger"
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	// Init user32.dll
	user32dll = syscall.NewLazyDLL("user32.dll")
	// For checking the state of a key
	procGetAsyncKeyState = user32dll.NewProc("GetAsyncKeyState")
	// To block user inputs
	procBlockInput = user32dll.NewProc("BlockInput")
	// To press keys
	procKeybdEvent = user32dll.NewProc("keybd_event")
	// To make pop up message boxes
	procMsgBox = user32dll.NewProc("MessageBoxW")
)

const (
	// https://github.com/micmonay/keybd_event/blob/master/keybd_windows.go
	_KEYEVENTF_KEYUP    = 0x0002
	_KEYEVENTF_SCANCODE = 0x0008
	VK_ALT              = 0x12 + 0xFFF

	// https://stackoverflow.com/a/71919136
	MB_SYSTEMMODAL = 0x00001000

	// https://github.com/winlabs/gowin32/blob/c9e40aa880584510eba1908f921a45317267ce9c/wrappers/winerror.go
	ERROR_SUCCESS syscall.Errno = 0
)

// WaitForKey() will wait for a key press and returns
// the UNIX timestamp of when the key was pressed
func WaitForKey(cfg config.Cfg) (timeStamp int64) {
	pressed := false

	for {
		// Iterate thru all keys
		for i := 0; i < 256; i++ {

			// Gets the state of a specific key
			// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getasynckeystate
			state, _, _ := procGetAsyncKeyState.Call(uintptr(i))

			// If the least significant bit is set, the key was pressed after the previous call to GetAsyncKeyState
			// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getasynckeystate#return-value
			if state&0x1 == 0 {
				continue
			}

			// Check if the key is in the ignoredKeys
			if intInSlice(i, cfg.IgnoredKeys) {
				// Key is in IgnoredKeys so ignore keystoke
				continue
			}

			timeStamp = time.Now().UnixMilli()
			pressed = true
		}

		// Better way of doing this? (wo/ named loops)
		if pressed {
			break
		}

		// To not have 100% CPU usage
		time.Sleep(time.Millisecond * 10)
	}

	// Return timestamp
	return timeStamp
}

// BlockInputFor() will block all user inputs for N seconds
func BlockInputFor(cfg config.Cfg) error {
	// Check if value is 0
	if cfg.BlockDuration == 0 {
		logger.LogWrite("Not blocking user inputs due to config")
	}

	// Block Inputs
	err := blockInput(true)
	if err != nil {
		return fmt.Errorf("blocking user inputs failed with: %s", err)
	} else {
		logger.LogWrite(fmt.Sprintf("Blocking user inputs for %d seconds", cfg.BlockDuration))
	}

	// Sleep for the specified duration
	time.Sleep(time.Duration(cfg.BlockDuration) * time.Second)

	// Unlock inputs
	err = blockInput(false)
	if err != nil {
		return fmt.Errorf("unblocking user inputs failed with: %s", err)
	} else {
		logger.LogWrite("No longer blocking user inputs")
	}

	return nil
}

// blockInput() will block the inputs via nf-winuser-blockinput
//
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-blockinput
func blockInput(block bool) error {
	// Block
	if block {
		r1, _, e1 := procBlockInput.Call(uintptr(int32(1)))
		// https://github.com/winlabs/gowin32/blob/master/wrappers/winuser.go
		if r1 == 0 {
			if e1 != ERROR_SUCCESS {
				return e1
			} else {
				// 536870951
				return syscall.EINVAL
			}
		}
		return nil

	} else {
		// Unblock
		r1, _, e1 := procBlockInput.Call(uintptr(int32(0)))
		if r1 == 0 {
			if e1 != ERROR_SUCCESS {
				return e1
			} else {
				return syscall.EINVAL
			}
		}
		return nil
	}
}

// SendKey() will press and release a key
func SendKey(key int) {
	downKey(key)
	time.Sleep(50 * time.Millisecond)
	upKey(key)
}

// downKey() will press a key
//
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-keybd_event
// https://github.com/micmonay/keybd_event/blob/master/keybd_windows.go
func downKey(key int) {
	flag := 0
	// Detect if the key code is virtual or no
	if key < 0xFFF {
		flag |= _KEYEVENTF_SCANCODE
	} else {
		key -= 0xFFF
	}
	vkey := key + 0x80
	procKeybdEvent.Call(
		uintptr(key),
		uintptr(vkey),
		uintptr(flag),
		0,
	)
}

// upKey() will release a key
//
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-keybd_event
// https://github.com/micmonay/keybd_event/blob/master/keybd_windows.go
func upKey(key int) {
	flag := _KEYEVENTF_KEYUP
	// Detect if the key code is virtual or no
	if key < 0xFFF {
		flag |= _KEYEVENTF_SCANCODE
	} else {
		key -= 0xFFF
	}
	vkey := key + 0x80
	procKeybdEvent.Call(
		uintptr(key),
		uintptr(vkey),
		uintptr(flag),
		0,
	)
}

// MessageBox() will pop open a Windows message box
//
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw
// https://stackoverflow.com/a/71919136
func MessageBox(cfg config.Cfg) {
	logger.LogWrite("Alerted via message box")

	titlePtr, _ := syscall.UTF16PtrFromString(cfg.AlertTitle)
	textPtr, _ := syscall.UTF16PtrFromString(
		fmt.Sprintf("%s %s",
			time.Now().Format(logger.TimeFormat),
			cfg.AlertMsg),
	)
	// normal
	flags := 0x00000000

	if cfg.AlertOnTop {
		// foreground / on top
		flags = MB_SYSTEMMODAL
	}

	syscall.SyscallN(
		procMsgBox.Addr(),
		0,
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(flags),
	)
}

// intInSlice will return true
// if a int is part of a []int
func intInSlice(value int, slice []int) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}

	return false
}
