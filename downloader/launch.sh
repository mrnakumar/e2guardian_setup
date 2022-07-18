go build main.go && ./main -userName="nkumar" \
-password="nkumar" \
-url="http://localhost:8080/eat" \
-downloadPath="/Users/narendra/workspace/e2guardian_setup/downloader/files" \
-headerEncryptKeyPath="/Users/narendra/workspace/e2guardian_setup/e2gserver/playground/headerPub" \
-bodyDecodeKeyPath="/Users/narendra/workspace/e2guardian_setup/syncer/playground/dataPri"
