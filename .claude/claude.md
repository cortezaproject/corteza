# Claude.md

The role of this file is to describe common mistakes and confusion points that agents might encounter as they work in this project.
If you ever encounter something in the project that surprises you, please alert the developer working with you and indicate that his is the case in the AgentMD file to help prevent future agent from having the same issue.

- If you use a new component make sure to register/import it. primevue-components for new PrimeVue ones. Otherwise in the file the component is used.

- Make sure FE respects the permissioning(RBAC) system.

- Always use translations for user facing text!
  If translations don't exist(i18n). First look in old-locale if they exist, use that (file/string whatever). If they are not in old-locale, make sure to add the new translations to locale.
