sudo cp /etc/e2guardian/examplef1.story /etc/e2guardian/group1.story
sudo mkdir /etc/e2guardian/lists/group1
sudo cp -r /etc/e2guardian/lists/example.group/* /etc/e2guardian/lists/group1/
sudo sed -i.bak 's/.*# comment out for production/#&/' /etc/e2guardian/e2guardianf1.conf
