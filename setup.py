import subprocess

# clean up ununsed dependencies
subprocess.run(["go", "mod", "tidy"])

# run main go file
subprocess.run(["go", "run", "."])
