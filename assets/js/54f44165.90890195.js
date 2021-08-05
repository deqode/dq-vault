"use strict";(self.webpackChunkdq_vault=self.webpackChunkdq_vault||[]).push([[152],{9170:function(t,e,a){a.r(e),a.d(e,{frontMatter:function(){return r},contentTitle:function(){return s},metadata:function(){return u},toc:function(){return d},default:function(){return h}});var n=a(7462),i=a(3366),o=(a(7294),a(3905)),l=["components"],r={sidebar_position:1},s=void 0,u={unversionedId:"getting-started/installation",id:"getting-started/installation",isDocsHomePage:!1,title:"installation",description:"This part of setting up dq-vault can be done using two methods. You may follow any one of your choices.",source:"@site/docs/getting-started/installation.md",sourceDirName:"getting-started",slug:"/getting-started/installation",permalink:"/docs/getting-started/installation",editUrl:"https://github.com/facebook/docusaurus/edit/master/website/docs/getting-started/installation.md",version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1},sidebar:"tutorialSidebar",previous:{title:"Introduction \ud83d\udc47\ufe0f",permalink:"/docs/intro"},next:{title:"configuration",permalink:"/docs/getting-started/configuration"}},d=[{value:"Vault installation",id:"vault-installation",children:[]},{value:"Get go files and Build plugin",id:"get-go-files-and-build-plugin",children:[]}],p={toc:d};function h(t){var e=t.components,a=(0,i.Z)(t,l);return(0,o.kt)("wrapper",(0,n.Z)({},p,a,{components:e,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"This part of setting up dq-vault can be done using two methods. You may follow any one of your choices."),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("h2",{parentName:"li",id:"method-1-"},"Method 1:-"),"Using ",(0,o.kt)("inlineCode",{parentName:"li"},"Docker")," to get your vault server up and running. You can find it in this ",(0,o.kt)("a",{parentName:"li",href:"https://github.com/deqode/dq-vault/tree/main/setup"},"link"),". We have provided the required docker files in the setup folder."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("h2",{parentName:"li",id:"method-2-"},"Method 2:-"),"Setting up Vault manually. The steps are given below in this document, starting from vault installation to creating your own vault server by using the CLI.")),(0,o.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,o.kt)("div",{parentName:"div",className:"admonition-heading"},(0,o.kt)("h5",{parentName:"div"},(0,o.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,o.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,o.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"info")),(0,o.kt)("div",{parentName:"div",className:"admonition-content"},(0,o.kt)("p",{parentName:"div"},"If you are already done with setting up the vault server using method 1, you may go directly to ",(0,o.kt)("strong",{parentName:"p"}," ",(0,o.kt)("a",{parentName:"strong",href:"https://deqode.github.io/dq-vault/docs/guides/usage"},"part 2"))," which elaborates the usage of the vault as an application server"))),(0,o.kt)("h2",{id:"vault-installation"},"Vault installation"),(0,o.kt)("p",null,"The first thing you need to do is to install vault to set-up a vault server."),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},"To install Vault, find the ",(0,o.kt)("a",{parentName:"p",href:"https://www.vaultproject.io/downloads.html"},"appropriate package")," for your system and download it. Vault is packaged as a zip archive.")),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},"After downloading Vault, unzip the package. Vault runs as a single binary named vault.")),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},"Copy the vault binary to your ",(0,o.kt)("inlineCode",{parentName:"p"},"PATH"),". In Ubuntu, PATH should be the ",(0,o.kt)("inlineCode",{parentName:"p"},"usr/bin")," directory.")),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},"To verify the installation, type vault in your terminal. You should see help output similar to the following:"),(0,o.kt)("pre",{parentName:"li"},(0,o.kt)("code",{parentName:"pre"},"  $ vault\n  Usage: vault <command> [args]\n\n  Common commands:\n      read        Read data and retrieves secrets\n      write       Write data, configuration, and secrets\n      delete      Delete secrets and configuration\n      list        List data or secrets\n      login       Authenticate locally\n      server      Start a Vault server\n      status      Print seal and HA status\n      unwrap      Unwrap a wrapped secret\n\n  Other commands:\n      audit          Interact with audit devices\n      auth           Interact with auth methods\n      lease          Interact with leases\n      operator       Perform operator-specific tasks\n      path-help      Retrieve API help for paths\n      policy         Interact with policies\n      secrets        Interact with secrets engines\n      ssh            Initiate an SSH session\n      token          Interact with tokens\n"))),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},"You can find the official installation guide ",(0,o.kt)("a",{parentName:"p",href:"https://www.vaultproject.io/intro/getting-started/install.html"},"here")))),(0,o.kt)("h2",{id:"get-go-files-and-build-plugin"},"Get go files and Build plugin"),(0,o.kt)("p",null,"Assuming that you have golang installed and your ",(0,o.kt)("inlineCode",{parentName:"p"},"GOPATH")," configured, get the plugin repository and run the build command in that folder:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-sh"},"  $ go build\n")),(0,o.kt)("p",null,"This will you give you a binary executable file with the name ",(0,o.kt)("inlineCode",{parentName:"p"},"dq-vault"),"."),(0,o.kt)("p",null,"Now move this binary file to a directory which the vault will use as its plugin directory. The plugin directory is where the vault looks up for available plugins."),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-sh"},"  $ mv dq-vault /etc/vault/plugins/vault_plugin\n")),(0,o.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,o.kt)("div",{parentName:"div",className:"admonition-heading"},(0,o.kt)("h5",{parentName:"div"},(0,o.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,o.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,o.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"info")),(0,o.kt)("div",{parentName:"div",className:"admonition-content"},(0,o.kt)("p",{parentName:"div"},"The above path is just an example, you can change the etc path to your own desired path."))))}h.isMDXComponent=!0}}]);