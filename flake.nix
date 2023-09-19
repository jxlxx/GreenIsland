{
  description = "A nix flake with a go dev environment";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system}; in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gotools
            pkgs.golangci-lint
            pkgs.gopls
            pkgs.go-outline
            pkgs.gopkgs
            pkgs.docker-compose
            pkgs.sqlc
            pkgs.postgresql
            pkgs.cue
            pkgs.terraform
            pkgs.tfswitch
          ];
        };
      }
  );
}
