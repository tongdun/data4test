<template>
  <div>
    <Table :columns="actionTableColumns" :data="actionTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'
import _ from 'lodash'
// import moment from 'moment'

@Component
export default class ActionList extends Vue {
  @Prop() isActions: boolean
  @Prop() isSending: boolean = false
  // @Prop() contentType: string
  @Prop() actionListData: Req.ActionListModel[]

  actionTableData: any[] = []
  actionType: string[] = ["sleep", "create_csv", "create_excel", "create_xlsx", "create_txt", "record_csv", "record_xls", "record_excel", "record_xlsx", "record_txt", "modify_file","create_hive_table_sql"]

  actionTableColumns: any[] = [
    {
      title: '',
      key: 'type',
      width: 200,
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.type,
                transfer: true,
                filterable: true,
                data: this.getOptions('type')
                // 'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'type', value)
                }
              }
            },
            this.actionType.map((item) => {
              return h('Option', {
                props: {value: item}
              }, item)
            })
        )
      }
    },
    {
      title: '',
      key: 'value',
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
              props: {
                value: params.row.value,
                transfer: true,
                data: this.getOptions("value")
                // 'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'value', value)
                }
              }
            }
        )
      }
    },
    {
      width: 80,
      renderHeader: (h: CreateElement, params: any) => {
        return h('Button', {
          props: {
            disabled: this.isSending,
            size: 'small',
            type: 'primary',
            icon: 'plus-round'
          },
          on: {
            click: () => {
              this.onAdd()
            }
          }
        })
      },
      render: (h: CreateElement, params: any) => {
        return h('Button', {
          props: {
            disabled: this.isSending,
            size: 'small',
            icon: 'close-round'
          },
          on: {
            click: () => {
              this.onDel(params.index)
            }
          }
        })
      }
    }
  ]

  typeOptions: string[] = []
  valueOptions: string[] = []

  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.actionTableColumns[0].title = this.$t('actionList.type')
    this.actionTableColumns[1].title = this.$t('actionList.value')
  }

  mounted() {
    this.syncTableData()
  }

  @Watch('actionListData')
  onListDataChange() {
    this.syncTableData()
  }

  // onFilterMethod(value: string, option: string) {
  //   return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  // }

  onAdd() {
    this.actionListData.push({
      type: '',
      value: '',
    })
    this.$emit('onChange')
    this.syncTableData()
  }

  onDel(index: number) {
    //删除指定的成员
    this.actionListData.splice(index, 1)
    this.$emit('onChange')

    this.syncTableData()
  }

  onChange(index: number, name: string, value: string) {
    if (name === 'type') {
      this.actionListData[index].type = value
    } else if (name === 'value') {
      this.actionListData[index].value = value
    }

    this.$emit('onChange')
  }

  getOptions(type: 'type' | 'value' ) {
    if (type==='type') {
      return this.typeOptions
    } else if (type==='value') {
      return this.valueOptions
    }

  }

  syncTableData() {
    //全量复制
    this.actionTableData = this.actionListData.map((m): any => {
      return {
        type: m.type,
        value: m.value
      }
    })
  }
}
</script>

