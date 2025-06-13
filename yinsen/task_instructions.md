# Yinsen Task Management Protocol
**Project: Mallon - Legal Document Automation System**

## 🚨 CRITICAL OPERATIONAL RULES

### Memory Management Protocol
- **Token Calculation**: Estimate task complexity before starting - break large tasks into subtasks
- **History Documentation**: Record detailed plan in `history.md` BEFORE beginning any task
- **Context Window Protection**: Stop and summarize to `history.md` when approaching memory limits
- **Session Continuity**: Update memories in `history.md` after EVERY completed task
- **Release Log**: Update release log in `releases.md` after EVERY completed task we can summarize features added, fixed, and changed to communmicate changes to the client and team

### Development Workflow Standards
1. **Feature Branch Creation**: ALWAYS create feature branch before starting work
2. **Version Control**: Increment version number on site after each task completion
3. **System Integrity**: Ensure `./start.sh`, `./restart.sh`, and `./stop.sh` remain functional
4. **Testing Protocol**: Start app with `./start.sh` after task completion for localhost testing
5. **QA Handoff**: Move ticket to `3_qa/` and alert for testing, then STOP and wait for response

## 📋 KANBAN WORKFLOW PROTOCOL

### Directory Structure Rules
```
1_queue/     - New tasks awaiting development
2_dev/       - Active development work  
3_qa/        - Ready for testing/review
4_done/      - Completed and verified tasks
```

### Task File Naming Convention
- **Defects**: `defect_{#}.md` 
- **Tasks**: `task_{#}.md`
- **Movement**: Files move through directories as work progresses

### Task Completion Tracking
- Mark completed tasks: `[x] TASK:{#} {DESC} - Completed: YYYY-MM-DD HH:MM`
- Only mark complete when task reaches `3_qa/` folder
- Update `history.md` with task completion details
- Update `releases.md` with task completion details only after the task is completed and verified and in the 4_done folder

## 🔄 TASK EXECUTION SEQUENCE

### Phase 1: Task Story Creation
1. **Review Queue**: Check existing files in `1_queue/`, `2_dev/`, `3_qa/` folders
2. **Generate Stories**: Create task files using `task_template.md` for missing tasks
3. **Priority Check**: ALWAYS check for defects first - defects take absolute priority
4. **Confirmation**: STOP and confirm understanding with task list and next action plan

### Phase 2: Development Execution  
1. **Defect Priority**: Complete ALL defects before starting feature tasks
2. **Task Selection**: Move highest priority task from `1_queue/` to `2_dev/`
3. **Implementation**: Follow task specifications with test-first approach when applicable  
4. **Quality Assurance**: Move completed work to `3_qa/` for review

### Phase 3: Validation & Completion
1. **Testing**: Start application and verify functionality
2. **Documentation**: Update `history.md` with completion details
3. **Version Management**: Increment version number
4. **Handoff**: Alert for testing and await approval before proceeding
5. **Completion**: Move completed work to `4_done/` and create PR after final verification
6. **Completion**: Update `releases.md` with task completion details only after the task is completed and verified and in the 4_done folder
7. **Completion**: create and push PR to github

## ⚠️ EXCLUSION RULES
- **IGNORE**: `4_done/archived_tasks/` folder - NOT part of active task management
- **PRESERVE**: Existing functionality - do not alter working features without explicit task requirements. Do not change UI elements or other pages not included in the current task. Maintain feature integrity and functionality, do not break existing features.

## 🎯 SUCCESS CRITERIA CHECKLIST
Before declaring any task complete, verify:
- [ ] Feature branch created and used
- [ ] Version number incremented on site
- [ ] Application starts successfully with `./start.sh`
- [ ] Task moved to `3_qa/` folder
- [ ] History updated with completion details
- [ ] No regression in existing functionality
- [ ] Ready for testing confirmation provided
- [ ] PR created and pushed to github
- [ ] Task moved to `4_done/` folder
- [ ] Release log updated with completion details
- [ ] Task moved to `5_archive/` folder

## 📊 TASK TEMPLATES

### Defect Entry Format
```
[ ] DEFECT:{#} {Description} - Priority: {HIGH/MEDIUM/LOW}
```

### Task Entry Format  
```
[ ] TASK:{#} {Description} - Estimated: {Size}
```

### Completion Entry Format
```
[x] TASK:{#} {Description} - Completed: YYYY-MM-DD HH:MM - Branch: feature/{task-name}
```

---

**Yinsen Protocol Version**: 2.0  
**Last Updated**: 2025-06-12  
**Compliance Required**: All development work must follow this protocol without exception