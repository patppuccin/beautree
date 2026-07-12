Output assets/beautree-demo.gif

Set FontFamily "JetBrainsMono Nerd Font"
Set FontSize 16
Set LetterSpacing 0
Set LineHeight 1.2
Set Width 1000
Set Height 800
Set Theme "GitHub Dark"
Set Padding 16
Set Framerate 30
Set Shell zsh
Set WindowBar Colorful
Set BorderRadius 16
Set Margin 16
Set MarginFill "#61AFEF"

Require beautree

Hide
Type "source ~/.zshrc"
Enter
Sleep 500ms
Type "SCAFFOLD=$(pwd)/scripts/generate-demo.sh"
Enter
Sleep 300ms
Type "bash $SCAFFOLD cleanup --dir ~/beautree-demo"
Enter
Sleep 500ms
Type "bash $SCAFFOLD create --dir ~/beautree-demo"
Enter
Sleep 500ms
Type "cd ~/beautree-demo"
Enter
Sleep 300ms
Type "clear"
Enter
Sleep 500ms
Show

Type "beautree --help"
Sleep 500ms
Enter
Sleep 3s
Type "clear"
Enter
Sleep 300ms

Type "beautree ."
Sleep 500ms
Enter
Sleep 3s
Type "clear"
Enter
Sleep 300ms

Type "beautree . -L 1"
Sleep 500ms
Enter
Sleep 3s
Type "clear"
Enter
Sleep 300ms

Type "beautree . -a"
Sleep 500ms
Enter
Sleep 3s
Type "clear"
Enter
Sleep 300ms

Type "beautree src/ --size"
Sleep 500ms
Enter
Sleep 3s

Hide
Type "bash $SCAFFOLD cleanup --dir ~/beautree-demo"
Enter
Sleep 500ms
Show