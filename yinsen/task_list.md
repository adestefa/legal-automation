Project: Mallon - Yinsen Core Memories


**SPECIAL INSTRUCTIONS**
- After all tasks, ./start.sh, restart.sh and ./stop.sh should still work.
- Warning, some tasks can be too large. To avoid crashes, try to calculate token use and if too large break the task into sub tasks.
- Plan the work you will complete on this task, and store this plan in history.md before you begin.
- Plan to stop a task and summarize to history.md before you run out of memory.
- After every task is complete update your memories in history.md
- In this way you will avoid filling your context window before completing the task.
- Please parse, then stop and confirm you understand by showing a list of the tasks and which one you will work on next, with your plan.



Yinsen, please find the below tasks we will work on togther. You will first help me create the full task stories from the list below.
1. Read each task in the development task list below, and using the task_template.md as a template create a story for each task in the list. Do the same for any defects listed, or skip them. Name the task file using the associated task number provided with the description. 
2. Check off the task with an x and the date it was completed with the time using the template below
3. Increment the version number shown on the site.
5. Only those completed should have an x and date. 
6. Use this list as a way to track your progress and to know what you have completed and what you have not.
7. Update the project history.md file with task number completed to reflect the changes you have made.
8. Before you begin a task, check for any defects. If there are defects, complete them first.
9. STOP. Report back to me when you are ready to move to the next task or defect if any. 

** Each defect should have an defect_#.md and each task should have a task_#.md file in the 1_queue folder.**
** as we comlpete each we move to 2_dev folder and then to 3_qa folder.**
** then we mark them with an x below when they reach 3_qa folder. **

NOTE: !!! IGNORE 4_done/archived_tasks/ folder !! it is not part of your task list. !!!

1. REMEMBER
1. Create a feature branch before any edits
2. Update the version number on the site
3. Run the ./start.sh to start the app so I can test it on localhost
4. If you need to stop the app run ./stop.sh
5. Move the ticket to 3_qa and alert me when you are ready for me to test. then stop and wait for my response.
6. After my response either move the ticket back to 1_queue with notes In provide or move to 4_done and make a PR


## Defect Tasks ##
1. [x] DEFECT:1 Navigation State Loss - User Selections Not Persisted During Back Navigation (Completed: 2025-06-13 10:15)
2. [x] DEFECT:1A Session Infrastructure - Add session management to Go handlers and state storage (GitHub Issue #2) (Completed: 2025-06-12 22:32)
3. [x] DEFECT:1B UI State Restoration - Update templates to read and display saved user selections (GitHub Issue #3) (Completed: 2025-06-12 23:10)
4. [x] DEFECT:1C Navigation Integration - HTMX navigation integration and end-to-end testing (GitHub Issue #4) (Completed: 2025-06-12 23:35)
5. [ ] DEFECT:2 Missing Content Tab Reports Errors When Data is Present (GitHub Issue #7)
6. [x] DEFECT:3 Step Icons Are No Longer Active (GitHub Issue #8) (Completed: 2025-06-13 10:45)


## Development Tasks
1. [ ] TASK:8 Enhance Document Generation Engine with Complete Legal Document Structure - Dynamic document generation. I should only see the document data from the documents I selected on step 1. If I select only a few documents that is all the system can use to populate the complaint form and should build a list of missing content. Ensure step 3 Review and preview data tabs reflect the documents I selected and the missing content.
2. [ ] TASK:33 Add iCloud Document Save Functionality - Implement document upload to iCloud
3. [x] TASK:34 replace "Restart Setup" button with a Back button that only moves the user to the last step. Every page should have a back button except the first and last pages. (Completed: 2025-06-12 22:11) 
4. [x] TASK:35 Return to use the user.json file to make it easier to admin the site and we can use it to add users and assign roles to them from future UI. (Completed: 2025-06-13 10:15) 
5. [x] TASK:36 Dynamic document preview based on selected documents - Document preview shows only user-selected documents with missing content analysis tab (Completed: 2025-06-13 00:30)