package cache

import "github.com/spf13/cobra"

var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Update and view info of the cache",
	Long: `opn uses a cache to speed up lookups of MIME types and paths to desktop files.
When desktop or mimeapps.list files are changed, either from the user manually changing it, or as
a result of the installation of a program, this cache can become out-of-date.

To update the cache, use "opn cache update".`,
}

func init() {
	CacheCmd.AddCommand(updateCacheCmd)
}
