#include <linux/module.h>
#define INCLUDE_VERMAGIC
#include <linux/build-salt.h>
#include <linux/elfnote-lto.h>
#include <linux/export-internal.h>
#include <linux/vermagic.h>
#include <linux/compiler.h>

BUILD_SALT;
BUILD_LTO_INFO;

MODULE_INFO(vermagic, VERMAGIC_STRING);
MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(".gnu.linkonce.this_module") = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};

#ifdef CONFIG_RETPOLINE
MODULE_INFO(retpoline, "Y");
#endif


static const struct modversion_info ____versions[]
__used __section("__versions") = {
	{ 0xbdfb6dbb, "__fentry__" },
	{ 0x336635eb, "single_open" },
	{ 0x5b8239ca, "__x86_return_thunk" },
	{ 0x40c7247c, "si_meminfo" },
	{ 0xbbbc1dd7, "seq_printf" },
	{ 0xd0da656b, "__stack_chk_fail" },
	{ 0x493a0ae2, "proc_create" },
	{ 0x92997ed8, "_printk" },
	{ 0x10a21d74, "remove_proc_entry" },
	{ 0x9a6c546e, "seq_read" },
	{ 0x1bdfd280, "module_layout" },
};

MODULE_INFO(depends, "");


MODULE_INFO(srcversion, "933B1C09F8ADDFC5BCDEB95");
