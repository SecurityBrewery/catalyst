import { Component, Mixins } from 'vue-property-decorator'
import {HorizontalBar, mixins} from 'vue-chartjs';
import ChartOptions from "chart.js";

@Component({
    extends: HorizontalBar,
    mixins: [mixins.reactiveProp],
    props: {
        chartOptions: {
            type: ChartOptions,
            default: null
        }
    },
})
export default class BarChart extends Mixins(mixins.reactiveProp, HorizontalBar) {
    mounted () {
        // @ts-expect-error chartOptions are not expected
        this.renderChart(this.chartData, this.chartOptions);
    }
}
