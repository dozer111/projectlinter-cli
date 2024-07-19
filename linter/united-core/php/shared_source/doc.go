package shared_source

// Resources shared by library and application
// During the evolution of the project, the following statement was formed
//
// The only common resources for application and library modes can be ONLY! configs for BumpLibrary and SubstituteLibrary rules
//
// Why?
// Because it would be quite strange that you avoid some libraries on application, while not avoiding them in libraries and vice versa
//
// Why, for example, the same .editorconfig for library/application cannot be put in this folder?
//
// Because most likely you take some application-tpl template as an example for application, respectively for library - library-tpl
// Accordingly, the references to the files(and maybe actually the files) for reconciliation in CVS will be different

// The full-filled examples:
// substitute => https://github.com/dozer111/projectlinter-core/blob/master/rules/dependency/substitute/full_example.yaml
// bump => https://github.com/dozer111/projectlinter-core/blob/master/rules/dependency/bump/full_example.yaml
