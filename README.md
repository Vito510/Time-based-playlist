# Time-based-playlist

Creates a playlist (m3u file) which ends exactly at a given time and immediately starts it using your default m3u playing application. Writen in golang + python (used for cacheing song durations), I would love to make it in C.

# How to use it
- run scan.py (it will probably require some more libraries so install them)
- input some paths that contain music and wait for it to complete
- run calc.go or calc.exe
- input end time [hour:minute]
- **[Bonus]** - you will probably first have to open the playlist manually so you can associate an app with the filetype 

