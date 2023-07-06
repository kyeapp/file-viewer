#!/bin/bash

rm file*
rm -r folder*

# Specify the directory where you want to create the folders and files
base_directory="."

# Loop to create three main folders
for ((i=1; i<=3; i++))
do
    # Create the main folder
    main_folder="$base_directory/folder$i"
    mkdir "$main_folder"
    touch "$base_directory/file$i.txt"

    # Loop to create three subfolders inside each main folder
    for ((j=1; j<=3; j++))
    do
        # Create the subfolder
        subfolder="$main_folder/subfolder$j"
        mkdir "$subfolder"
        touch "$main_folder/sub_file$j.txt"

        # Loop to create three files inside each subfolder
        for ((k=1; k<=3; k++))
        do
            # Create the file
            touch "$subfolder/sub_sub_file$k.txt"
        done
    done
done
