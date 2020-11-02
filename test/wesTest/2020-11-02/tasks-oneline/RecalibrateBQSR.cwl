cwlVersion: v1.0

class: CommandLineTool
id: sort_bam
requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement 
    ramMin: 2000

inputs:
  input_bam: 
    type: File
  indexed_reference_fasta: 
    type: File 
    secondaryFiles: [.64.amb, .64.ann, .64.bwt, .64.pac, .64.sa, .64.alt, .dict, .fai]
  indexed_reference_dpsnp_hg38:
    type: File
  output_recalibrate_bqsr:
    type: string 
  
outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*.recal.table'

arguments:
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      -R $(inputs.indexed_reference_fasta) --use-original-qualities -I $(inputs.input_bam)  --known-sites $(inputs.indexed_reference_dpsnp_hg38) -O $(input.output_recalibrate_bqsr)

baseCommand: [gatk, BaseRecalibrator]
