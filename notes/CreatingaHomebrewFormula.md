# Creating a Homebrew Formula


Notes on building a Homebrew Formula for `atmotool`


Create a release, tar.gz it up, get the shasum.

	tar -czf atmotool-1.6.0.tar.gz atmotool
	shasum -a 256 atmotool-1.6.0.tar.gz


Update the [homebrew-akana](https://github.com/ghchinoy/homebrew-akana) repo's `atmotool.rb` Formula to include the shasum and the release.

Upload the tar.gz to the [releases](https://github.com/ghchinoy/atmotool/releases) of [atmotool](https://github.com/ghchinoy/atmotool).
Update the `atmotool.rb` with the URL to the OS X release.


# Output

	brew tap akana/atmotool git@bitbucket.org:akana/homebrew-atmotool.git
	brew install atmotool


## Create a public Git repo

Bitbucket git@bitbucket.org:akana/homebrew-atmotool.git

Contains
* Homebrew Formula
* download of atmotool

## Tar

	tar -cvzf atmotool-1.1.1.tar.gz atmotool-1.1.1

## Modify the Formula

### SHA1

	shasum -a 1 atmotool-1.1.1.tar.gz



## Testing

Remove the tap

	brew untap akana/atmotool


# References

* [brew tap](https://github.com/Homebrew/homebrew/blob/master/share/doc/homebrew/brew-tap.md) help file 
* [direnv Formula](https://github.com/Homebrew/homebrew/blob/master/Library/Formula/direnv.rb)
* [gpm Formula](https://github.com/Homebrew/homebrew/blob/master/Library/Formula/gpm.rb) - [gpm](https://github.com/pote/gpm)
* [kong tap Formula](https://github.com/Mashape/homebrew-kong/blob/master/Formula/kong.rb)