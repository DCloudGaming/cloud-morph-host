#Region ;**** Directives created by AutoIt3Wrapper_GUI ****
#AutoIt3Wrapper_Compile_Both=y
#AutoIt3Wrapper_UseX64=y
#EndRegion ;**** Directives created by AutoIt3Wrapper_GUI ****
Func writeToBackgroundNotepad()
    Run("notepad.exe")

    ; Wait 10 seconds for the Notepad window to appear then minimize
    Local $hWnd = WinWait("[CLASS:Notepad]", "", 10)
	WinSetState ("[CLASS:Notepad]", "", @SW_MINIMIZE)

    ; Wait for 2 seconds.
    Sleep(2000)

    ; Send a string of text to the edit control of Notepad. The handle returned by WinWait is used for the "title" parameter of ControlSend.
    ControlSend($hWnd, "", "Edit1", "DCloud gaming")

    ; Wait for 2 seconds.
    Sleep(2000)

	ControlClick($hWnd, "", "Edit1", "left", 1, 0, 0)

    Sleep(2000)

    ; Send a string of text to the edit control of Notepad. The handle returned by WinWait is used for the "title" parameter of ControlSend.
    ControlSend($hWnd, "", "Edit1", "DCloud")

	MsgBox(0, "Hello World!", "Done")
EndFunc


writeToBackgroundNotepad()