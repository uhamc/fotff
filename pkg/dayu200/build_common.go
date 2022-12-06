package dayu200

const (
	preCompileCMD = `rm -rf prebuilts/clang/ohos/darwin-x86_64/clang-480513;rm -rf prebuilts/clang/ohos/windows-x86_64/clang-480513;rm -rf prebuilts/clang/ohos/linux-x86_64/clang-480513;bash build/prebuilts_download.sh`
	compileCMD    = `echo 'start' && export NO_DEVTOOL=1 && export CCACHE_LOG_SUFFIX="dayu200-arm32" && export CCACHE_NOHASHDIR="true" && export CCACHE_SLOPPINESS="include_file_ctime" && ./build.sh --product-name rk3568 --ccache --build-target make_all --build-target make_test --gn-args enable_notice_collection=false`
)

var imgList = []string{
	"out/rk3568/packages/phone/images/MiniLoaderAll.bin",
	"out/rk3568/packages/phone/images/boot_linux.img",
	"out/rk3568/packages/phone/images/parameter.txt",
	"out/rk3568/packages/phone/images/system.img",
	"out/rk3568/packages/phone/images/uboot.img",
	"out/rk3568/packages/phone/images/userdata.img",
	"out/rk3568/packages/phone/images/vendor.img",
	"out/rk3568/packages/phone/images/resource.img",
	"out/rk3568/packages/phone/images/config.cfg",
	"out/rk3568/packages/phone/images/ramdisk.img",
	"out/rk3568/packages/phone/images/chipset.img",
	"out/rk3568/packages/phone/images/sys_prod.img",
	"out/rk3568/packages/phone/images/chip_prod.img",
	"out/rk3568/packages/phone/images/updater.img",
	"out/rk3568/packages/phone/updater/bin/updater_binary",
}
