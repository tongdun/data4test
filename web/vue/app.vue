<template>
    <div class="layout">
        <Layout>
            <Content :style="{margin:'0px 0px 0','z-index':0}">
                <div class="layout-content">
                    <div class="layout-content-main">
                        <router-view></router-view>
                    </div>
                </div>
            </Content>
        </Layout>
    </div>
</template>

<style>
.layout {
    border: 1px solid #d7dde4;
    background: #f5f7f9;
    position: relative;
    border-radius: 4px;
    overflow: hidden;
}

.layout-content {
    min-height: 200px;
    margin: 15px;
    overflow: hidden;
    background: #fff;
    border-radius: 4px;
}

.layout-content-main {
    padding: 10px;
}
</style>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'

@Component
export default class App extends Vue {
    mounted() {
    }

}
</script>
