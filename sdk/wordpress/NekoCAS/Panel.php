<?php defined( 'ABSPATH' ) or exit; ?>
<div class="wrap">
    <h2>NekoCAS</h2>
    <hr>
</div>

<div>
    <form method="POST" action="">
        <table>
            <tr>
                <td>
                    <label for="login_text">后台登录按钮文本</label>
                </td>
                <td>
                    <input type="text" class="regular-text" name="login_text"
                           value="<?php _e( $value['login_text'] ); ?>"/>
                </td>
            </tr>

            <tr>
                <td>
                    <label for="cas_domain">NekoCAS 地址</label>
                </td>
                <td>
                    <input type="text" class="regular-text" name="cas_domain"
                           value="<?php _e( $value['cas_domain'] ); ?>">
                </td>
            </tr>

            <tr>
                <td>
                    <label for="secret">Secret</label>
                </td>
                <td>
                    <input type="text" class="regular-text" name="secret"
                           value="<?php _e( $value['secret'] ); ?>">
                </td>
            </tr>
        </table>
        <br>
        <input type="submit" name="submit" value="修改配置" class="button button-primary"/>
    </form>
</div>
