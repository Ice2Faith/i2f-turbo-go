If WScript.Arguments.Count < 1 Then
    WScript.Echo "usage: cscript wps_to_pptx.vbs <inputPath> [outputPath]"
    WScript.Quit 1
End If

Dim fso
Set fso = CreateObject("Scripting.FileSystemObject")

Dim inputPath, outputPath
inputPath = WScript.Arguments(0)

If WScript.Arguments.Count >= 2 Then
    outputPath = fso.GetAbsolutePathName(WScript.Arguments(1))
Else
    outputPath = inputPath & ".pptx"
End If

inputPath = fso.GetAbsolutePathName(inputPath)
outputPath = fso.GetAbsolutePathName(outputPath)

Dim app, doc

On Error Resume Next

Set app = CreateObject("KWPP.Application") ' invoke WPS ppt
app.Visible = True
app.WindowState = 2   ' 2 means min window

Set doc = app.Presentations.Open(inputPath)
' FileFormat = 24 means pptx format
doc.SaveAs outputPath, 24 

If Not doc Is Nothing Then
    doc.Close
End If
If Not app Is Nothing Then
    app.Quit
End If

Set doc = Nothing
Set app = Nothing