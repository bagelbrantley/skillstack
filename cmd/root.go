package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	serverURL  string
	jsonOutput bool
)

var rootCmd = &cobra.Command{
	Use:           "skillstack",
	Short:         "skillstack - CLI for SkillStack Guided Projects",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		jsonOutput = viper.GetBool("json")
		if s := viper.GetString("server"); s != "" {
			serverURL = s
		}
	},
	Long: style.BoxStyle.Render(
		style.TitleStyle.Render("SkillStack CLI") + "\n" +
			"Develop, test, and submit Guided Projects directly on your local workspace.",
	),
}

func ExecuteContext(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/skillstack/credentials.json)")
	rootCmd.PersistentFlags().StringVarP(&serverURL, "server", "S", "https://ais-dev-rz3jr65n6cnsvt3dwpu4pl-368667717597.us-east1.run.app", "SkillStack server URL")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")

	_ = viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(home, ".config", "skillstack"))
			viper.SetConfigType("yaml")
			viper.SetConfigName(".skillstack")
		}
	}

	viper.SetEnvPrefix("SKILLSTACK")
	viper.AutomaticEnv()

	viper.SetDefault("server", "https://ais-dev-rz3jr65n6cnsvt3dwpu4pl-368667717597.us-east1.run.app")

	_ = viper.ReadInConfig()
}

func GetServerURL() string {
	if serverURL != "" {
		return serverURL
	}
	if s := viper.GetString("server"); s != "" {
		return s
	}
	return "https://ais-dev-rz3jr65n6cnsvt3dwpu4pl-368667717597.us-east1.run.app"
}

func IsJSONMode() bool {
	return jsonOutput || viper.GetBool("json")
}
