sudo cp /etc/e2guardian/examplef1.story /etc/e2guardian/group1.story
dir="/etc/e2guardian/lists/group1"
sudo mkdir "${dir}" || { echo "Failed to make directory ${dir}. Exiting."; exit 1; }
sudo cp -r /etc/e2guardian/lists/example.group/* "${dir}/"
sudo sed -i.bak 's/.*# comment out for production/#&/' /etc/e2guardian/e2guardianf1.conf

# TODO: Allow school online education websites and enable blanket block.
