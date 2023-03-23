go build calc.go
pyinstaller --upx-dir C:\Users\Vito\Downloads\upx-4.0.2-win64 -F scan.py
mv calc.exe dist/play.exe
mv dist/scan.exe dist/scan.exe
cp template/config.json dist/config.json
rm build -r