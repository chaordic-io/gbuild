{ pkgs ? import <nixpkgs> {
  overlays = [ (self: super: {
    # nodejs = super.nodejs-10_x;
    # jre = super.jdk11;
    # xcodebuild = super.tmux;
  }) ];
} }:

pkgs.mkShell {
  name="dev-environment";
  buildInputs = [
   pkgs.go
   pkgs.oh-my-zsh
   pkgs.vscode
   pkgs.cloc
  ];
  shellHook = ''
    echo "Welcome to your dev env"
    export PATH=/usr/bin:$PATH
    '';
    }
    # need to add original path first to not override xcrun etc