# Code generation package

It consumes YAMLs from all known locations and provides them to
templates

 - determinate all known locations for yaml files
 - determinate all known output types (go?, html?, adocs?)
   - can we read this from the filename?
 - define structure for all locations
 - define outputs
   - can be a static file
   - can have placeholders that are assembled from the structure, and yaml file contents
