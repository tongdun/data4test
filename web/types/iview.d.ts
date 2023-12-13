import Vue, { VNode } from 'vue'

declare module "vue/types/vue" {
    interface Vue {
        $Modal: iv.Modal
        $Message: iv.Message
    }
}

declare namespace iv {
    interface Modal {
        info: (config: ModalConfig) => void
        success: (config: ModalConfig) => void
        warning: (config: ModalConfig) => void
        error: (config: ModalConfig) => void
        confirm: (config: ModalConfig) => void
        remove: () => void
    }

    interface ModalConfig {
        title?: string
        content?: string
        render?: (this: Vue, h: any) => VNode
        width?: number
        okText?: string
        cancelText?: string
        loading?: boolean
        scrollable?: boolean
        onOk?: () => void
        onCancel?: () => void
    }

    interface Message {
        info: (config: MessageConfig) => void
        success: (config: MessageConfig) => void
        warning: (config: MessageConfig) => void
        error: (config: MessageConfig) => void
        loading: (config: MessageConfig) => void
        config: (config: MessageConfig) => void
        destory: () => void
    }

    interface MessageConfig {
        content?: string
        duration?: number
        onClose?: () => void
        closable?: boolean
    }
}