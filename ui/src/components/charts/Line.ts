import { Component, Mixins } from 'vue-property-decorator'
import {Line, mixins} from 'vue-chartjs';
import ChartOptions from "chart.js";

@Component({
    extends: Line,
    mixins: [mixins.reactiveProp],
    props: {
        chartOptions: {
            type: ChartOptions,
            default: null,
        }
    },
})
export default class LineChart extends Mixins(mixins.reactiveProp, Line) {
    mounted () {
        // @ts-expect-error chartOptions are not expected
        this.renderChart(this.chartData, this.chartOptions);
    }
}
