from glob import glob
from rich.progress import track
from os.path import isdir
import json
import audioread
paths = ["C:\\Users\\Vito\\Music", "D:\\Music"]
audio = ["flac","mp3","wav","ogg","m4a","opus","aac"]

data = []
types = {}

items = []
for p in paths:
    items.extend(glob(p+"\\**", recursive=True))

for i in track(items):
    if isdir(i) or i.split(".")[-1] not in audio:
        continue

    t = i.split(".")[-1]

    if t not in types:
        types[t] = 1
    else:
        types[t] += 1

    with audioread.audio_open(i) as f:
        data.append(
            {
                "path": i,
                "duration": f.duration
            }
        )

data = sorted(data, key= lambda x: x["duration"])


with open("songs.json", "w", encoding='utf-8') as f:
    json.dump(data, f, indent=4)


print(types)