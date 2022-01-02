###
# This file can be used to install e2guardian.
###

# Create user
echo "Adding user e2guardian..."
useradd -m e2guardian

echo "Installing e2guardian..."
./install.sh

# Download e2guardian
wget https://github.com/e2guardian/e2guardian/archive/refs/heads/v5.4.zip || exit

cd ./e2guardian-5.4

# Build e2guardian
echo "Compiling e2guardian..."
sudo ./autogen.sh || exit
sudo ./configure '--prefix=/usr' '--enable-clamd=yes' '--with-proxyuser=e2guardian' '--with-proxygroup=e2guardian' '--sysconfdir=/etc' '--localstatedir=/var' '--enable-icap=yes' '--enable-commandline=yes' '--enable-email=yes' '--enable-ntlm=yes' '--mandir=${prefix}/share/man' '--infodir=${prefix}/share/info' '--enable-pcre=yes' '--enable-sslmitm=yes' 'CPPFLAGS=-mno-sse2 -g -O2' || exit

sudo make || exit
sudo make install || exit

# Copy files for group1
./group1_setup.sh



## Operations stuff
# Configure log rotate
./add_logrotate.sh

# Make e2guardian run on boot
sudo cp ./data/scripts/e2guardian.service /etc/systemd/system/
sudo systemctl enable e2guardian

# Configure iptables to make sure access is only via e2guardian
./iptables.sh
