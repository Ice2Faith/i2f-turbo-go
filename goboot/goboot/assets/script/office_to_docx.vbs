If WScript.Arguments.Count < 1 Then
    WScript.Echo "usage: cscript office_to_docx.vbs <inputPath> [outputPath]"
    WScript.Quit 1
End If

Dim fso
Set fso = CreateObject("Scripting.FileSystemObject")

Dim inputPath, outputPath
inputPath = WScript.Arguments(0)

If WScript.Arguments.Count >= 2 Then
    outputPath = fso.GetAbsolutePathName(WScript.Arguments(1))
Else
    outputPath = inputPath & ".docx"
End If

inputPath = fso.GetAbsolutePathName(inputPath)
outputPath = fso.GetAbsolutePathName(outputPath)

Dim app, doc

On Error Resume Next

Set app = CreateObject("Word.Application") ' invoke Office word
app.Visible = False
app.WindowState = 2   ' 2 means min window
app.DisplayAlerts = False

Set doc = app.Documents.Open(inputPath)
' FileFormat = 12 means docx format
doc.SaveAs outputPath, 12

If Not doc Is Nothing Then
    doc.Close False
End If
If Not app Is Nothing Then
    app.Quit
End If

Set doc = Nothing
Set app = Nothing