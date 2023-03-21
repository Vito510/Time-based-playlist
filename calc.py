import json
from datetime import timedelta
import subprocess
from datetime import datetime
from numpy.random import shuffle
from numpy.random import choice
from os.path import abspath
from os import getcwd
import audioread

def get_duration(path: str):
    with audioread.audio_open(path) as f:
        return f.duration


MAIN_PATH = abspath(getcwd())

startupinfo = subprocess.STARTUPINFO()
startupinfo.dwFlags |= subprocess.STARTF_USESHOWWINDOW

def subset_sum(numbers, target, partial=[]):
    s = 0
    for item in partial:
        s += item["duration"]

    # check if the partial sum is equals to target
    if s == target: 
        return partial
    if s >= target:
        return  # if we reach the number why bother to continue
    
    for i in range(len(numbers)):
        n = numbers[i]
        remaining = numbers[i+1:]
        r = subset_sum(remaining, target, partial + [n]) 
        if r is not None:
            return r
   
n = datetime.now()
target_time = datetime(n.year, n.month, n.day, 13, 0, 0)
# target = timedelta(minutes=20, seconds=20).total_seconds()
target = (target_time - datetime.now()).total_seconds()

if target < 0:
    print("Error: target_time can't be in the past")
    exit()


MIN_LENGTH = timedelta(minutes=2).total_seconds()
MAX_LENGTH = timedelta(minutes=5).total_seconds()

FORCE_LAST_SONG = ["C:\\Users\\Vito\\Music\\Outer Wilds OST (.FLAC)\\10. End Times.flac"]

if len(FORCE_LAST_SONG) > 0:
    LAST_SONG = choice(FORCE_LAST_SONG)
    d = get_duration(LAST_SONG)
    if d > target:
        print("Error: couldn't fit last song")
        LAST_SONG = None
    else:
        target -= d
        LAST_SONG = {
            "path": LAST_SONG,
            "duration": d
        }
        print(LAST_SONG)
else:
    LAST_SONG = None

# Precision
target = round(target,1)

data = json.loads(open('songs.json', 'r', encoding='utf-8').read())

# Remove all songs longer than time target
temp = []
for song in data:
    if song["duration"] < MIN_LENGTH or song["duration"] > MAX_LENGTH:
        continue
    elif song["duration"] <= target:
        temp.append(song)
    else:
        break
data = temp

shuffle(data)


playlist = subset_sum(data, target)

# Add forced last song
if LAST_SONG is not None:
    playlist.append(LAST_SONG)


with open("playlist.m3u", "w", encoding='utf-8') as f:
    f.write("#EXTM3U\n")
    for song in playlist:
        f.write(song["path"]+"\n")



subprocess.Popen(["C:\Program Files (x86)\MusicBee\MusicBee.exe",f'{MAIN_PATH}\\playlist.m3u'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, startupinfo=startupinfo, creationflags=subprocess.CREATE_NO_WINDOW)

print("Created playlist")
print(f"\tEnd time: {target_time}\n\tLength: {str(timedelta(seconds=target)).split('.')[0]}\n\tMin song length: {timedelta(seconds=MIN_LENGTH)}\n\tMax song length: {timedelta(seconds=MAX_LENGTH)}")
print(f"\n\tItems:")
for item in playlist:
    print('\t\t'+item["path"].split("\\")[-1])
