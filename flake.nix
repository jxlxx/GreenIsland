{
  description = "A nix flake with a go dev environment";

  inputs = {
    unstable.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self,  unstable, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = unstable.legacyPackages.${system}; in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.gotools
            pkgs.golangci-lint
            pkgs.gopls
            pkgs.go-outline
            pkgs.gopkgs
            pkgs.go_1_21
          ];
        };
      }
  );
}
