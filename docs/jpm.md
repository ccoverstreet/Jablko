# Jablko Package Manager (JPM)

JPM is the package manager for users wanting to install Jablko Modules into their Jablko system. The main method of module distribution is through GitHub repositories, where JPM downloads a zip version of the target tag to users Jablko install, installs any dependencies, and updates their configuration file.

## Contents
- (Installing Packages)[#installing-packages]
- (Uninstalling Packages)[#uninstalling-packages]
- (Reinstalling Packages)[#reinstalling-packages]

## Installing Packages

To install packages there are two available command syntaxes. The easier of the two is made for GitHub repositories only: `./jpm install author/repo tag target_name`. For example, `./jpm install ccoverstreet/Jablko-Interface-Status v1.1.0 interface_status` will install the Jablko Interface Status module Version 1.1.0 to interface_status. JPM will automatically add an entry in the "jablko_config.json" file which contains the module information and default config. 

*ATTENTION*: Some modules may need their config values in jablko_config.json to modified or set or else they will prevent Jablko from starting. For example, the ccoverstreet/Jablko-Mantle-RGB needs the user to change the default `controller_ip` value from `null` to a string containing the IP address of the RGB light controller (ex. `192.168.1.101`).

## Uninstalling Packages

To remove packages from Jablko, use the syntax `./jpm uninstall module_name` where the module name is the installed name of the module (NOT the original repository or source name). This command will remove the module from your "jablko_config.json". If you wish to reinstall a module without deleting configuration use (Reinstall)[#reinstalling-packages].

## Reinstalling Packages

To reinstall a package, use the syntax `./jpm reinstall module_name`. This will delete the current module's directory and redownload the module from the `repo_archive` specified in the "jablko_config.json" file.
