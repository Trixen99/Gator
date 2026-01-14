Gator - Blog Aggregator 

Requirements:
Go
PostgreSQL


Installation (Linux):

Go:

(please note this software was created on Arch Linux limited support for any other distro or operating system)

First you will need to install Go, to test whether you have this installed already type this into your terminal
go version

You should see something like this appear in the terminal:
go version go1.25.5 linux/amd64

If nothing appears you will need to install Go.

For Arch Linux this can be done with:
sudo pacman -Syu
sudo pacman -S go

Other wise for debian try these two commands:
sudo apt update
sudo apt install golang-go

Once finished make sure that when you type "go version" you get a version number of at least "Go 1.21".


PostgreSQL:

Now you will need to install PostgreSQL, to test if you already have this software type the following into your terminal;
psql --version

You should get something like;
psql (PostgreSQL) 18.1

If nothing appears then you still need to install it.

For arch linux systems you can use the following command to install PostgreSQL;
sudo pacman -S postgresql

Otherwise for debian based Operating Systems use;
sudo apt update
sudo apt install postgresql postgresql-contrib

Once you think it is installed type "psql --version" and make sure that the version number is at least 15"

Finally you will now need to start PostgreSQL in the background.

On arch Linux this can be done with;
sudo systemctl start postgresql
If you want to make sure it starts everytime you boot your computer also use this command;
sudo systemctl enable postgresql

On Debian based distro's this can be done with;
sudo service postgresql start


You will now need to create an password for PostgreSQL using this command;
sudo passwd postgres 

once this command in entered into the console you will then be asked for a password.

You will then need to type these commands to finish the setup of PostgreSQL (copy and paste one at a time);
sudo -u postgres psql

CREATE DATABASE gator;

\c gator

Finally type the command below (if you do not wish to use 'postgres' as your password change it now before entering the command but make sure you can remember it)
ALTER USER postgres PASSWORD 'postgres';

you can then type "exit" to get out of the shell.



Downloading and setting up Gator:


Now you can create a local clone of the repsitory and build the files. This can be done by creating a new folder where you would like to download the files. Navigate to that folder in your terminal and then typing the following commands;
"git clone https://github.com/Trixen99/Gator.git"
git build

You will now be able to interact with the Blog Aggregator however we haven't finished setting it up yet so is likely to error. 

The commands you need to run are these;

nano ~/.gatorconfig.json

you will then need to paste this line of text into the file;
{
  "db_url": "postgres://<password>:postgres@localhost:5432/gator",
  "current_user_name": "username_goes_here"
}

replace '<password>' with the password I told you to remember previously and then save by using "Ctrl + x" on your keyboard, followed by typing "Y" and pressing "Enter"


You should now be able to start using Gator. 

Commands:

Here are a few commands you can run

Users (Displays all users in the program);
./gator users

Register (allows the creation of a new user and logs you into that user, Important - Name must be one word);
./gator register <insert user_name>

Addfeed (create a new feed for the currently logged in user, feed url must link to an xml, Iportant - Feed_Name must be one word or wrapped in '' marks)
./gator addfeed <insert feed_name> <inser feed_url>

Feeds (Displays all users in the program)
./gator feeds

Follow (allows a user to follow a feed not created by them)
./gator follow <insert feed_url>

Following (displays all feeds the currently logged in user is following)
./gator following

Unfollow (remove a feed from the currently logged in user)
./gator unfollow <insert feed_url>

agg (start the aggregator, this is an infinite loop and will attempt to call the urls you provided in the feed connected to your current users, do not set the interval time to low here or you maybe considered a hacker/malicious actor for the dos attack you are performing, format for time is <number><unit initial> (for example 30s or 5m). you may need to open another terminal window as the one that runs this command will be locked until you use "Ctrl + C);
./gator agg <insert time interval>

browse (will allow you to browse your downloaded posts, can give this command a number to limit the amount of post you are given (defaults to 2). Posts are always oredered in order of creation (newest first));
./gator browse <optional limit>
