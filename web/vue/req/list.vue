<template>
    <div>
        <Table :columns="runTableColumns" :data="runTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')" @on-selection-change="onSelect"></Table>
    </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'
import _, {hasIn} from 'lodash'
// import moment from 'moment'

@Component
export default class RunList extends Vue {
    @Prop() isHeaders: boolean
    @Prop() isSending: boolean = false
    @Prop() contentType: string
    @Prop() runListData: Req.RunListModel[]

    multiTestValue: string[] = []
    nameOptions: string[] = []
    valueTypeOptions: string[] = []
    testValueOptions: string[] = []
    egValueOptions: string[] = []
    descOptions: string[] = []
    isDisableOptions: boolean[] = []
    paramTypeOptions: string[] = ["string", "boolean", "integer", "object", "array"]
    isMustOptions: string[] = ["yes", "no"]
    inValueTypeOptions: string[] = ["{self}", "{Time}", "{Date}", "{Timestamp(-3)}", "{Timestamp(1)}", "{Int(0,10)}", "{Int(16,120)}", "{Str(8)}", "{Str(64)}", "{Rune(8)}","{Rune(64)}", "{Name}", "{BankNo}", "{Address}", "{Email}", "{IDNo}", "{Mobile}", "{Province}", "{City}", "{DeviceType}", "{QQ}", "{Sex}", "{Age}", "{Diploma}", "{IntStr10}", "{Int}", "{Result}", "{Income}", "{YorN}", "{Level}", "{DayBegin}","{MonthBegin}","{YearBegin}","{DayEnd}","{MonthEnd}","{YearEnd}","{DayBegin(-10)}","{DayEnd(-7)}"]

    runTableData: any[] = []
    runTableColumns: any[] = [
        {
            type: 'selection',
            width: 60,
            align: 'center'
        },
        {
            title: '',
            key: 'name',
            width: 150,
            render: (h: CreateElement, params: any) => {
                return h('AutoComplete', {
                    props: {
                        // disabled: this.getDisabled(params.index),
                        value: params.row.name,
                        transfer: true,
                        data: this.getOptions('name'),
                        'filter-method': this.onFilterMethod
                    },
                    on: {
                        'on-change': (value: string) => {
                            this.onChange(params.index, 'name', value)
                        }
                    }
                })
            }
        },
        {
        title: '',
        width: 100,
        key: 'valueType',
        render: (h: CreateElement, params: any) => {
          return h('AutoComplete', {
                props: {
                  value: params.row.valueType,
                  transfer: true,
                  data: this.getOptions('valueType'),
                  'filter-method': this.onFilterMethod
                },
                on: {
                  'on-change': (value: string) => {
                    this.onChange(params.index, 'valueType', value)
                  }
                }
              },
              this.paramTypeOptions.map((item) => {
                return h('Option', {
                  props: {value: item}
                }, item)
              })
          )
        }
      },
        {
        title: '',
        key: 'isMust',
        width: 80,
        render: (h: CreateElement, params: any) => {
          return h('AutoComplete', {
                props: {
                  value: params.row.isMust,
                  transfer: true,
                  data: this.getOptions("isMust"),
                  'filter-method': this.onFilterMethod
                },
                on: {
                  'on-change': (value: string) => {
                    this.onChange(params.index, 'isMust', value)
                  }
                },
              },
              this.isMustOptions.map((item) => {
                return h('Option', {
                  props: {value: item}
                }, item)
              }))
        }
      },
        {
            title: '',
            key: 'testValue',
            // width: 160,
            render: (h: CreateElement, params: any) => {
                return h('Select', {
                    props: {
                        value: params.row.testValue,
                        transfer: true,
                        multiple: true,
                        data: this.getOptions("testValue"),
                        filterable: true,
                        'allow-create': true,
                        'on-create': this.onAddInValueType(params.row.testValue),
                    },
                    on: {
                      'on-change': (value: string) => {
                        this.onChange(params.index, 'testValue', value)
                      }
              }
                },
                    this.inValueTypeOptions.map((item) => {
                      return h('Option', {
                        props: {label: item, value: item}
                      }, item)
                    }))
            }
        },
        {
        title: '',
        key: 'egValue',
        width: 120,
        render: (h: CreateElement, params: any) => {
          return h('Input', {
            props: {
              // disabled: this.getDisabled(params.index),
              value: params.row.egValue
            },
            on: {
              'on-change': (value: string) => {
                this.onChange(params.index, 'egValue', value)
              }
            }
          })
        }
      },
        {
        title: '',
        key: 'desc',
        width: 120,
        render: (h: CreateElement, params: any) => {
          return h('Input', {
            props: {
              // disabled: this.getDisabled(params.index),
              value: params.row.desc
            },
            on: {
              'on-change': (value: string) => {
                this.onChange(params.index, 'desc', value)
              }
            }
          })
        }
      },
        {
            width: 60,
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
    created() {
        this.onLocale()
    }

    @Watch('$i18n.locale')
    onLocale() {
        //注意 tableColumns 的处理位置需要与上面初始化定义一致
        this.runTableColumns[1].title = this.$t('list.name')
        this.runTableColumns[2].title = this.$t('list.valueType')
        this.runTableColumns[3].title = this.$t('list.isMust')
        this.runTableColumns[4].title = this.$t('list.testValue')
        this.runTableColumns[5].title = this.$t('list.egValue')
        this.runTableColumns[6].title = this.$t('list.desc')
    }

    mounted() {
        this.syncRunTableData()
        this.onContentTypeChange()
        // this.onRunListDataChange()

        if (this.isHeaders) {
            this.nameOptions = ['Content-Type', 'User-Agent']
            this.valueTypeOptions = [
                'text/plain',
                'text/html',
                'application/json',
                'application/javascript',
                'application/xml',
                'application/x-www-form-urlencoded'
            ]
        }
    }

    @Watch('runListData')
    onRunListDataChange() {
        this.syncRunTableData()
    }

    @Watch('contentType')
    onContentTypeChange() {
        if (!this.isHeaders) {
            return
        }

        if (this.contentType == '') {
            return
        }

        if (this.contentType == 'text') {
            //text状态下，不加任何content-type
            for (let i = 0; i < this.runListData.length; i++) {
                if (this.runListData[i].testValue.toLowerCase() === 'content-type') {
                    this.onDel(i)
                    break
                }
            }

            return
        }

        let isFound = false

        _.forEach(this.runListData, v => {
            if (v.name.toLowerCase() === 'content-type') {
                v.testValue = this.contentType
                isFound = true
                return
            }
        })

        if (!isFound) {
            this.runListData.push({
                isDisable: false,
                name: '',
                egValue: '',
                testValue: '',
                isMust: 'no',
                desc: '',
                valueType: ''
            })
        }

        this.syncRunTableData()
    }

    onAddInValueType(val) {
      _.forEach(val, v => {
        if (this.inValueTypeOptions.indexOf(v)<0) {
          this.inValueTypeOptions.push(v)
        }
      })
  }
  setDefaultOptions(val) {
      this.inValueTypeOptions.push(val)
  }

    onFilterMethod(value: string, option: string) {
        return option.toLowerCase().indexOf(value.toLowerCase()) != -1
    }

    onSelect(selection: Req.RunListModel[]) {
        _.forEach(this.runListData, p => {
            p.isDisable = true

            let s = _.find(selection, function(s) {
                return p.name === s.name
            })
            if (s) {
                p.isDisable = false
            }
        })
        this.$emit('onChange')
    }

    onAdd() {
        this.runListData.push({
            isDisable: false,
            name: '',
            valueType: 'string',
            testValue: '',
            isMust: 'no',
            egValue: '',
            desc: ''
        })
        this.$emit('onChange')
        this.syncRunTableData()
    }

    onDel(index: number) {
        //删除指定的成员
        this.runListData.splice(index, 1)
        this.$emit('onChange')

        this.syncRunTableData()
    }

    onChange(index: number, name: string, value: string) {
        if (name === 'name') {
            this.runListData[index].name = value
        } else if (name === 'valueType') {
          this.runListData[index].valueType = value
        }  else if (name === 'isMust') {
            this.runListData[index].isMust = value
        } else if (name === 'testValue') {
            this.runListData[index].testValue = value
        } else if (name === 'egValue') {
          this.runListData[index].egValue = value
        } else if (name === 'desc') {
          this.runListData[index].desc = value
        }
      this.$emit('onChange')
    }

    getOptions(type: 'name' | 'valueType' | 'isMust' |'testValue'|  'egValue' | 'desc') {
      if (type==='name') {
        return this.nameOptions
      }  else if (type ==='valueType') {
        return this.valueTypeOptions
      } else if (type ==='isMust') {
        return this.isMustOptions
      } else if (type ==='testValue') {
        return this.testValueOptions
      } else if (type ==='egValue') {
        return this.egValueOptions
      } else if (type ==='desc') {
        return this.descOptions
      }

    }

    getDisabled(index: number) {
        if (this.isSending) {
            return true
        }
        return this.runListData[index].isDisable
    }

    syncRunTableData() {
        //全量复制
        this.runTableData = this.runListData.map((m): any => {
            return {
                isDisable: m.isDisable,
                name: m.name,
                valueType: m.valueType,
                isMust: m.isMust,
                testValue: m.testValue,
                egValue: m.egValue,
                desc: m.desc,
                _checked: !m.isDisable,
            }
        })
    }

}
</script>

