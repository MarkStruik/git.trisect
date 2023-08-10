# Git Trisect

the better bisect

## Goals

- [ ] Ask for a folder
- [ ] Ask for a last valid date estimate
- [ ] Display a list of check-in's for that date and let the user pick the start point
- [ ] check if there is a package.json 
  - [ ] ask user if they want to run a command after the switch
  - [ ] if package.json dependencies changed 
    - [ ] ask user if they want to npm ci
- [ ] wait for confirmation ( good/bad/skip/close )
- [ ] continue until find the wrong version or close