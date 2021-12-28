while IFS='\n' read -r package; do
   sudo apt-get install $package -y || continue
done < requirements
