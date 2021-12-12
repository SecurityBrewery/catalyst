import { Component, Mixins } from "vue-property-decorator";
import { mixins, Pie } from "vue-chartjs";
import ChartOptions from "chart.js";

@Component({
    extends: Pie,
    mixins: [mixins.reactiveProp],
    props: {
        chartOptions: {
            type: ChartOptions,
            default: null
        }
    }
})
export default class DoughnutChart extends Mixins(mixins.reactiveProp, Pie) {
    mounted() {
        // @ts-expect-error chartOptions are not expected
        this.renderChart(this.chartData, this.chartOptions);
    }
}
