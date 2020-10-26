Run go get github.com/urfave/cli 
https://github.com/urfave/cli
Make sure the package is downloaded either in the GOPATH or in GOROOT

Run go build main.go
Run cp main /usr/local/bin to make the cli tool globally available
main -h  shows help
main -v shows version
main getRequest -u <url>
main doProfile -u <url> -p <profile-times>
    
   
   
Tried some other sites
    
shreya90@Shreyas-MacBook-Air cloudfare % main dP -u http://google.com/imAs -p 3
Total number of times the url is requested for Profiling : 3 
Maxinum time for requesting: 1038 
Mininum time for requesting: 1013 
Mean time for requesting: 1021 
Maximum response length in bytes: 1720 
Minimum response length in bytes: 1720 
Percentage of Success: 0 
Error codes for failures: [   404 404 404]


shreya90@Shreyas-MacBook-Air cloudfare % main dP -u http://www.google.com/imghp -p 3                             
Total number of times the url is requested for Profiling : 3 
Maxinum time for requesting: 1031 
Mininum time for requesting: 1016 
Mean time for requesting: 1021 
Maximum response length in bytes: 49121 
Minimum response length in bytes: 49100 
Percentage of Success: 100 
Error codes for failures: [  ] 

