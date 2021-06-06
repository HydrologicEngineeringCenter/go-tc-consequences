# go-tc-consequences
A command line application to support computation of consequences from tropical cyclones and hurricanes.

# consequences
The consequences compute is based on depth alone and uses the default USACE damage functions for the 40 occupancy types consistent with HEC-FIA and HEC-LifeSim.

# Running the program.
1. Download the project from github.
2. Install Docker Desktop
3. Launch Docker Desktop
4. Go to settings -> Resources -> File Sharing and share the directory where the code from github is stored.
5. Download Visual Studio Code
6. Launch Visual Studio Code
7. Install the extension Remote Development
8. Open the go-tc-consequence project with visual studio
9. click on the green "><" icon in the lower left corner of Visual Studio Code
10. Select the option Reopen This workspace in Container
11. Select the container "Docker.dev"
12. Once the container is built and the project is reloaded install the extension for go, when it asks you in the lower right corner if you want delve, say yes to all, once they have installed, open a terminal and type go build main.go
13. Once the binary main is built type ./main -h to get help on the command line arguments
14. Stage any data you wish to compute with in the ./data/ folder
15. To update the code to get latest from github you need to be sure you are connected to the repo. type git remote -v to determine any connections
16. If no remotes are detected type git remote add origin https://github.com/HydrologicEngineeringCenter/go-tc-consequences.git
17. To pull latest type git pull origin main.
18. After pulling, rebuild the binary with go build main.go
