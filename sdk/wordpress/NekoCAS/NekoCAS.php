<?php
/*
Plugin Name: NekoCAS
Description: 统一登录插件
Version: 1.0.0
Author: NekoWheel Works
*/

defined( 'ABSPATH' ) or exit;
define( 'PLUGIN_URL', plugin_dir_url( __FILE__ ) );

// 初始化数据
add_option(
	'nekocas_data',
	array(
		'login_text' => '使用 Neko 账号登录',
		'cas_domain' => '',
		'secret'     => ''
	),
	null,
	'yes'
);

// 左侧显示插件菜单
add_action( 'admin_menu', 'nekocas_admin_menu' );
function nekocas_admin_menu() {
	add_menu_page(
		__( 'NekoCAS' ),     //展开菜单名
		__( 'NekoCAS' ),     //主菜单名
		'administrator',   //权限
		'nekocas_admin',
		'nekocas_admin_page',     //显示界面
		'dashicons-admin-links',      //图标
		30      //位置
	);
}

// 管理页面
function nekocas_admin_page() {
	if ( isset( $_POST['login_text'], $_POST['cas_domain'] ) ) {
		update_option(
			'nekocas_data',
			array(
				'login_text' => $_POST['login_text'],
				'cas_domain' => $_POST['cas_domain'],
				'secret'     => $_POST['secret']
			)
		);
	}
	$value = get_option( 'nekocas_data' );
	require_once( 'Panel.php' );
}

// 修改登录页面
function nekocas_login_form() {
	$config      = get_option( 'nekocas_data' );
	$callbackURL = $config['cas_domain'] . 'login?service=' . urlencode( site_url() . '/login/cas' );
	echo '<a href="' . htmlentities( $callbackURL ) .
	     '"><input type="button" class="button button-primary" style="width: 100%;" value="' .
	     htmlentities( $config['login_text'] ) .
	     '"></a><br><br><br>';
}

add_action( 'login_form', 'nekocas_login_form' );

// 回调路由
function nekocas_rewrite_basic() {
	add_rewrite_rule( 'login/cas$', 'index.php', 'top' );
}

function nekocas_flush_rewrites() {
	nekocas_rewrite_basic();
	flush_rewrite_rules();
}

register_deactivation_hook( __FILE__, 'flush_rewrite_rules' );
register_activation_hook( __FILE__, 'nekocas_flush_rewrites' );

add_action( 'init', 'nekocas_init' );
function nekocas_init() {
	$config = get_option( 'nekocas_data' );
	$query  = http_build_query( array( 'service' => $config['secret'], 'ticket' => $_GET['ticket'] ) );

	$ch = curl_init();
	curl_setopt( $ch, CURLOPT_URL, $config['cas_domain'] . 'validate?' . $query );
	curl_setopt( $ch, CURLOPT_RETURNTRANSFER, 1 );
	$output = curl_exec( $ch );
	curl_close( $ch );
	$response = json_decode( $output, true );

	if ( is_null( $response ) || ! $response['success'] ) {
		return;
	}

	$email = $response['email'];
	$user  = get_user_by_email( $email );
	if ( $user === false ) {
		flush();
		die( '电子邮箱不存在！' );
	}
	wp_set_auth_cookie( $user->ID, false, '' );
	do_action( 'wp_login', $user->user_login, $user );
	wp_redirect( admin_url() );
	exit;
}