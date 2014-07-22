SHELL = /bin/bash
DESTDIR = ./dist
EXT = $(DESTDIR)/Pinboard.lbext
PB = $(DESTDIR)/Pinboard.lbaction
PBB = $(DESTDIR)/Pinboard\ Browse.lbaction
PBCONTENTS = $(PB)/Contents
PBBCONTENTS = $(PBB)/Contents

all: clean bundle refresh

ext: all
	@echo "Making .lbext"
	@install -d $(EXT)/Contents/Resources/Actions
	@install -pm 0644 ./src/Info.plist $(EXT)/Contents
	@cp -r $(PB) $(EXT)/Contents/Resources/Actions
	@cp -r $(PBB) $(EXT)/Contents/Resources/Actions
	@ditto -ck --keepParent $(EXT) $(EXT).zip
	@$(RM) -r $(EXT)
	@mv -f $(EXT).zip $(EXT)

clean:
	@$(RM) -r ./dist

bundle: golb gopb PB PBB

pb: golb gopb PB refresh

pbb: golb gopb PBB refresh

golb:
	go install github.com/nbjahan/go-launchbar

gopb:
	go install github.com/nbjahan/go-pinboard

PB:
	@echo "Creating the Pinboard.lbaction"
	@install -d ${PBCONTENTS}/{Resources,Scripts}
	@install -pm 0644 ./src/pinboard/Info.plist $(PBCONTENTS)
	@install -pm 0644 ./resources/* $(PBCONTENTS)/Resources
	go build -o $(PBCONTENTS)/Scripts/pinboard ./src/pinboard

PBB:
	@echo "Creating the Pinboard Browse.lbaction"
	@install -d $(PBBCONTENTS)/{Resources,Scripts}
	@install -pm 0644 ./src/pinboard-browse/Info.plist $(PBBCONTENTS)
	@install -pm 0644 ./resources/* $(PBBCONTENTS)/Resources
	go build -o $(PBBCONTENTS)/Scripts/pinboard ./src/pinboard-browse

edit:
	@subl -n ./ ~/code/go/src/github.com/nbjahan/go-launchbar/ ~/code/go/src/github.com/nbjahan/go-pinboard/

refresh:
	@echo "Refreshing the LaunchBar"
	@osascript -e 'run script "tell application \"LaunchBar\" \n repeat with rule in indexing rules \n if name of rule is \"Actions\" then \n update rule \n exit repeat \n end if \n end repeat \n activate \n end tell"'

.PHONY: all clean bundle pb pbb golb gopb edit refresh