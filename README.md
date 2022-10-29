# [pixiv](https://www.pixiv.net/)-crawler  

## start crawler with command line arguments
```  
NAME:
   image downloader - download image from pixiv 

USAGE:
   main.exe [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -d value, --download value    input IllustID to download
   -u value, --url value         input pixiv url to download
   -a value, --author value      author id (default: 0)
   --user value, --userid value  input user id (default: 0)
   -n value, --name value        author name
   -f, --following               following
   -r, --recommend               recommend illust
   -s, --stars                   download stars
   --rk, --ranking               ranking illust
   --help, -h                    show help
   --version, -v                 print the version


```

## 封装 api来自项目 https://github.com/everpcpc/pixiv

