cwlVersion: v1.0

class: CommandLineTool 
id: align_bwa_mem 

requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement
    ramMin: 4096
    coresMin: 4

inputs:
  input_fastqs1:
    type: File[]
  input_fastqs2:
    type: File[]
  indexed_reference_fasta: 
    type: File 
    secondaryFiles: [.64.amb, .64.ann, .64.bwt, .64.pac, .64.sa, .64.alt, .dict, .fai]
  output_aligned_reads:
    type: string

outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*.bam'
      
arguments: 
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      $(inputs.indexed_reference_fasta) $(inputs.input_fastqs1) $(inputs.input_fastqs2) | samtools view -Shb -o $(inputs.output_aligned_reads) -
baseCommand: [bwa, mem]
