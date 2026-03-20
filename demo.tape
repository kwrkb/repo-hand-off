Output demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 960
Set Height 540
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 90ms

Type "clear"
Enter
Sleep 500ms

Type "echo '# repo-hand-off: capture and share your dev state'"
Enter
Sleep 1s

Type "./handoff export"
Enter
Sleep 2s

Type "head -30 HANDOFF.md"
Enter
Sleep 3s

Type "./handoff prompt --format xml | head -20"
Enter
Sleep 3s

Type "./handoff diff"
Enter
Sleep 3s

Type "echo '# pipe directly to your AI agent'"
Enter
Sleep 500ms
Type "./handoff prompt | pbcopy"
Enter
Sleep 2s
