#!/bin/bash


## create a new post in your site with a random post content
name=$(uuidgen)
content=$(fortune)

hugo new $name.md
echo $content >> /sitePath/post/$name.md

## generate the static content of the website
hugo --baseUrl="http://remsh.github.io/"


## git commit changes and push it to upstream repo
cd /sitePath/public
#git init
#git remote add origin https://github.com/remsh/remsh.github.io.git
git add .
git commit -m "Add Post $name"
git push -u origin master

