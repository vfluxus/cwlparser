cwlVersion: v1.0

class: CommandLineTool 
id: merge_bam_alignment 

requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement
    ramMin: 4096
    coresMin: 4

inputs:
  aligned_bam:
    type: File
  unaligned_bam:
    type: File
  indexed_reference_fasta:
    type: File
    secondaryFiles: [.64.amb, .64.ann, .64.bwt, .64.pac, .64.sa, .64.alt, .dict, .fai]
  output_merged_bam:
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
      ALIGNED=$(inputs.aligned_bam) UNMAPPED=$(inputs.unaligned_bam) O=$(inputs.output_merged_bam) R=$(inputs.indexed_reference_fasta)
baseCommand: [picard, MergeBamAlignment]