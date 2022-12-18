/*
 * Copyright (c) 2022 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dayu200

const (
	preCompileCMD = `rm -rf prebuilts/clang/ohos/darwin-x86_64/clang-480513;rm -rf prebuilts/clang/ohos/windows-x86_64/clang-480513;rm -rf prebuilts/clang/ohos/linux-x86_64/clang-480513;bash build/prebuilts_download.sh`
	compileCMD    = `echo 'start' && export NO_DEVTOOL=1 && export CCACHE_LOG_SUFFIX="dayu200-arm32" && export CCACHE_NOHASHDIR="true" && export CCACHE_SLOPPINESS="include_file_ctime" && ./build.sh --product-name rk3568 --ccache --build-target make_all --gn-args enable_notice_collection=false`
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
	// "out/rk3568/packages/phone/images/chipset.img",
	"out/rk3568/packages/phone/images/sys_prod.img",
	"out/rk3568/packages/phone/images/chip_prod.img",
	"out/rk3568/packages/phone/images/updater.img",
	// "out/rk3568/packages/phone/updater/bin/updater_binary",
}
