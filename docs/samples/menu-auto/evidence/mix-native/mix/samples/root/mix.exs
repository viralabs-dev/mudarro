defmodule Sample.MixProject do
  use Mix.Project
  def project, do: [app: :mudarro_root, version: "0.1.0", elixir: "~> 1.20", deps: []]
  def application, do: [extra_applications: [:logger]]
end
