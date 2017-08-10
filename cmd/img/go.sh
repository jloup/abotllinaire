#!/bin/sh

./img -i poem.txt

wkhtmltoimage --quality 100 poem.html poem.png
