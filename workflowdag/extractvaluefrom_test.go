package workflowdag

import (
	"fmt"
	"testing"
)

func TestExtractValueFrom(t *testing.T) {
	ExtractValueFrom("SM=$(inputs.output_basename) RG=$(inputs.output_basename)  LIBRARY_NAME=lib1 PLATFORM=ILLUMINA F1=$(inputs.input_fastq1) F2=$(inputs.input_fastq2) O=$(inputs.output_fastq2bam_bam)")
	fmt.Println("----")
	ExtractValueFrom("$(inputs.input_fastq1) $(inputs.input_fastq2) -o .")
	fmt.Println("----")
	ExtractValueFrom(" -t 6 $(inputs.indexed_reference_fasta) $(inputs.input_fastq1) $(inputs.input_fastq2) | samtools view -Shb -o $(inputs.output_aligned_reads) ")
	fmt.Println("----")
	ExtractValueFrom("-I $(inputs.input_bam) -bqsr $(inputs.input_recalibrate_bqsr) -O $(inputs.output_apply_bqsr)")
	fmt.Println("----")
	ExtractValueFrom("--num_shards=6 --model_type=WES --ref=$(inputs.indexed_reference_fasta) --reads=$(inputs.input_bam) --output_vcf=$(inputs.output_vcf_dv) --output_gvcf=$(inputs.output_gvcf_dv)")
}
